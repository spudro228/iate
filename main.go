package main

import (
	"context"
	"encoding/json"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
	"golang.org/x/sync/errgroup"
	"iate/ai"
	"iate/models"
	"iate/product"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// region http_handlers
func testHandler(openaiClient *openai.Client, service *product.InMemoryProductService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
		w.WriteHeader(http.StatusOK)
	}
}
func productHandlerUpdate(service *product.InMemoryProductService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var productObj product.Product
		err = json.Unmarshal(body, &productObj)
		if err != nil {
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}

		err = service.UpdateProduct(productObj)
		if err != nil {
			http.Error(w, "Error update product", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func productsHandlerGetAll(service *product.InMemoryProductService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		products, err := service.GetAllProducts()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Print(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(products)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Print(err)
			return
		}
	}
}

func productsHandlerDelete(service *product.InMemoryProductService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var deleteModel models.DeleteProductById
		err = json.Unmarshal(body, &deleteModel)
		if err != nil {
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}

		err = service.DeleteProduct(deleteModel.Guid)
		if err != nil {
			http.Error(w, "Error deleting product", http.StatusBadRequest)
			return
		}
	}
}

//endregion

// region bot_handler
func botMessageGlobalHandler(vk *api.VK, handlerCallback func(string) error) func(context.Context, events.MessageNewObject) {
	return func(ctx context.Context, object events.MessageNewObject) {
		msgText := object.Message.Text
		log.Print(msgText)

		//пинпонг для тестирования
		if msgText == "hello" {
			b := params.NewMessagesSendBuilder()
			b.Message("hello!")
			b.RandomID(0)
			b.PeerID(object.Message.PeerID)

			r, err := vk.MessagesSend(b.Params)
			if err != nil {
				log.Print(err, r)
			}

			return
		}

		err := handlerCallback(msgText)
		var responseMsgTxt string
		if err != nil {
			responseMsgTxt = "Что-то не получилось, попробуйте еще раз позже или измените запрос."
		} else {
			responseMsgTxt = "Сохранил!"
		}
		b := params.NewMessagesSendBuilder()
		b.Message(responseMsgTxt)
		b.RandomID(0)
		b.PeerID(object.Message.PeerID)

		_, err = vk.MessagesSend(b.Params)
		if err != nil {
			log.Print(err)
		}
	}
}

//endregion

// region ai_interceptor
func tryToParseAndSaveInfoFromUser(
	ctx context.Context,
	openaiClient *openai.Client,
	service *product.InMemoryProductService,
	request string) error {
	//Я съел 1000 грамм риса с 100 граммами кабачков и 300 грамм куриной грудки.
	content, err := ai.Send(ctx, openaiClient, request)

	if err != nil {
		log.Printf("ChatCompletion error: %v\n", err)
		return nil //todo: отвечать реплаем что не сохранил сообщение
	}

	var products []product.Product
	err = json.Unmarshal([]byte(content), &products)
	if err != nil {
		return err
	}

	timeNow := time.Now()
	for _, productObj := range products {
		productObj.Guid = uuid.NewString()
		productObj.CreatedAt = timeNow
		err = service.SaveProduct(productObj)
		if err != nil {
			log.Printf("Can't save productObj %+v\n", productObj)
		} else {
			log.Printf("Save productObj %+v\n", productObj)
		}
	}

	return nil
}

// endregion

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	VKToken, existsVk := os.LookupEnv("VK_TOKEN")
	if !existsVk {
		log.Print("VK_TOKEN empty. Exit.")
		return
	}

	OpenAiToken, existsOAi := os.LookupEnv("OPEN_AI_TOKEN")

	if !existsOAi {
		log.Print("OPEN_AI_TOKEN empty. Exit.")
		return
	}

	ctxBackground := context.Background()
	//грасефул шатдаунт не работает, используй сигкилл
	ctx, stop := signal.NotifyContext(ctxBackground, os.Interrupt, os.Kill)
	defer stop()

	g, ctx := errgroup.WithContext(ctx)
	productService := product.NewInMemoryProductService()
	client := openai.NewClient(OpenAiToken)

	//region http_server
	g.Go(func() error {
		http.HandleFunc("/products/update", productHandlerUpdate(productService))
		http.HandleFunc("/products/getAll", productsHandlerGetAll(productService))
		http.HandleFunc("/products/delete", productsHandlerDelete(productService))
		http.HandleFunc("/test", testHandler(client, productService))

		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
			return err
		}

		return nil
	})
	//endregion

	//region vk_bot_server
	vk := api.NewVK(VKToken)

	g.Go(func() error {
		lp, err := longpoll.NewLongPollCommunity(vk)
		if err != nil {
			panic(err)
		}

		lp.Goroutine(true) // для обработки множеста сообщений, с тандартной либе уже строено оборачивание обработчиков в горутины

		lp.MessageNew(
			botMessageGlobalHandler(
				vk,
				func(msg string) error {
					return tryToParseAndSaveInfoFromUser(
						ctx,
						client,
						productService,
						msg,
					)
				},
			),
		)
		return lp.RunWithContext(ctx)
	})
	//endregion

	err := g.Wait()
	if err != nil {
		log.Fatal(err)
	}

}
