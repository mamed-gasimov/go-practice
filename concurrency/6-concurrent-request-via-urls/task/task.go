package main

func main() {
	urls := []string{
		"https://google.com",
		"https://yandex.ru",
		"https://amazon.com",
		"https://youtube.com",
	}

	process(urls)
}

// реализовать паралельные запросы по адресам из списка
// подсчитать количество для каждого StatusCode ответа
// предусмотреть возможность отмены запроса по таймауту
func process(urls []string) {}
