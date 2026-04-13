package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Backend представляет биржу с именем и задержкой ответа,
// имитирующей реальное время обработки запроса.
type Backend struct {
	Name  string
	Delay time.Duration
}

// DoRequest симулирует HTTP-запрос к бирже.
// Если контекст отменён раньше, чем истечёт задержка — возвращает пустую строку.
func (w *Backend) DoRequest(ctx context.Context) string {
	select {
	case <-ctx.Done():
		// Контекст завершён (таймаут или явная отмена) — прерываем запрос
		return ""
	case <-time.After(w.Delay):
		// Задержка истекла — запрос завершён успешно
	}

	return fmt.Sprintf("Response of backend %s", w.Name)
}

func main() {
	// Три биржи с разными временами ответа: NYSE — медленнее всего, MOEX — быстрее всего
	backends := []*Backend{
		{Name: "NYSE", Delay: 4 * time.Second},
		{Name: "MOEX", Delay: 1 * time.Second},
		{Name: "LSE", Delay: 2 * time.Second},
	}

	// Создаём контекст с таймаутом 3 секунды.
	// NYSE (4s) не успеет ответить — её запрос будет прерван.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Читаем все результаты из канала — получим только ответы,
	// успевшие прийти до истечения таймаута (MOEX и LSE).
	var results []string
	for name := range getFirstName(ctx, backends) {
		results = append(results, name)
	}

	// Ждём завершения контекста перед выводом, чтобы убедиться,
	// что все горутины завершили работу.
	<-ctx.Done()
	fmt.Println(results)
}

// getFirstName отправляет запросы ко всем бэкендам параллельно
// и возвращает канал, в который горутины пишут результаты по мере готовности.
// Канал закрывается, когда все горутины завершатся.
func getFirstName(ctx context.Context, backends []*Backend) <-chan string {
	names := make(chan string)
	var wg sync.WaitGroup

	for _, backend := range backends {
		wg.Add(1)
		go func() {
			// Отправляем запрос; если контекст истёк — DoRequest вернёт ""
			names <- backend.DoRequest(ctx)
			wg.Done()
		}()
	}

	// Отдельная горутина закрывает канал после завершения всех запросов,
	// чтобы цикл range в main мог корректно выйти.
	go func() {
		wg.Wait()
		close(names)
	}()

	return names
}
