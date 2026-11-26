# ArtiVana · Filter Snippet

Модуль фильтрации и отображения каталога товаров, реализованный на **Go** (backend) и **Nuxt 4** + **Pinia** (frontend). Проект демонстрирует разработку гибкой системы фильтрации товаров, работу с PostgreSQL и организацию серверной логики в отдельных модулях.

## 🎯 Цель

Создать гибкий и переиспользуемый модуль фильтрации для e-commerce приложений, который:
- поддерживает фильтрацию товаров по категориям, подкатегориям, брендам, характеристикам, цене и рейтингу;
- формирует метаданные фильтров с количеством доступных товаров для каждого значения;
- учитывает взаимосвязи фильтров и актуализирует доступные варианты в зависимости от выбранных параметров;
- обеспечивает сортировку и пагинацию каталога;
- предоставляет фронтенду структурированные данные для построения интерактивного интерфейса.

## ⚙️ Технические особенности

- **Backend**: Go, Gin, PostgreSQL, pgx — обработка HTTP-запросов, работа с базой данных и реализация логики фильтрации;
- **Frontend**: Nuxt 4 (Composition API), Pinia, модульные компоненты и тесты для каждой ключевой части интерфейса;
- **Testing**: Unit-тесты frontend части для логики фильтров, компонентов и утилит.
- **Data Layer**: отдельные репозитории для получения товаров, категорий, брендов и характеристик; работа с динамическими параметрами фильтрации.
- **Filter Metadata**: формирование фасетов с количеством доступных товаров, диапазонами цен и характеристиками.
- **API Documentation**: Swagger — документация эндпоинтов, параметров запросов и структур ответов.

---

> Код выложен исключительно в ознакомительных целях.


## 📁 Структура репозитория

### 🧩 Frontend

| Файл | Описание |
|------|----------|
| `components/card/BaseInput.vue` | Инпут для карточки товара: ввод/вывод количества и добавление в корзину |
| `components/card/BaseInput.test.ts` | Тесты логики BaseInput.vue |
| `components/card/CardVertical.vue` | Карточка товара в вертикальном виде |
| `components/card/CardVertical.test.ts` | Тесты логики CardVertical.vue |
| `components/card/ImagesSlider.vue` | Слайдер изображений товара для карточке |
| `components/card/Tooltip.vue` | Копирование артикула товара с подсказкой |
| `components/common/AnimateNumber.vue` | Анимация числовых значений (например, цен) |
| `components/common/AnimateNumber.test.ts` | Тесты логики AnimateNumber.vue |
| `components/common/BaseInput.vue` | Универсальный кастомный input |
| `components/common/BaseInput.test.ts` | Тесты логики BaseInput.vue |
| `components/common/OptionSelector.vue` | Кастомный checkbox/radio-компонент |
| `components/common/RangeSlider.vue` | Кастомный range-слайдер для выбора диапазона |
| `components/common/RangeSlider.test.ts` | Тесты логики RangeSlider.vue |
| `components/filter/Filters.vue` | Основной компонент фильтров товаров |
| `components/filter/Filters.test.ts` | Тесты логики Filters.vue |
| `components/filter/Filtration.vue` | Обёртка фильтров для десктопной версии |
| `components/filter/Heading.vue` | Верхний контейнер фильтра |
| `components/filter/Heading.test.ts` | Тесты логики Heading.vue |
| `components/filter/List.vue` | Список карточек товаров |
| `components/filter/Main.vue` | Главная обёртка модуля фильтрации |
| `components/filter/Pagination.vue` | Компонент пагинации для фильтрации |
| `components/filter/RangeInput.vue` | Двойной input для ручного ввода ценового диапазона |
| `components/filter/RangeInput.test.ts` | Тесты логики RangeInput.vue |
| `components/filter/SelectedFilters.vue` | Отображение активных (выбранных) фильтров |
| `components/filter/SelectedFilters.test.ts` | Тесты логики SelectedFilters.vue |
| `pages/catalog/[subcategory]/index.vue` | Страница фильтрации товаров по подкатегории |
| `pages/catalog/[subcategory]/[subsubcategory]/index.vue` | Страница фильтрации товаров по подподкатегории |
| `stores/filter.ts` | Хранилище Pinia для управления состоянием фильтра |
| `stores/filter.test.ts` | Тесты логики filter.ts |
| `typescript/interfaces.ts` | Интерфейсы для моделей данных фильтра, товаров и компонентов |
| `typescript/types.ts` | Типы TypeScript для проекта |
| `utils/price-format.ts` | Утилита форматирования цен с разделением пробелами |
| `utils/price-format.test.ts` | Тесты логики price-format.ts |
| `utils/word-ending.ts` | Утилита для склонения окончаний слов |

---

### 🧱 Backend

| Файл | Описание |
|------|----------|
| `internal/modules/categories/model.go` | Модели категорий, подкатегорий и подподкатегорий |
| `internal/modules/categories/repository.go` | Интерфейс и реализация доступа к данным категорий |
| `internal/modules/dictionaries/model.go` | Модели брендов, характеристик и их значений |
| `internal/modules/dictionaries/repository.go` | Интерфейс и реализация получения справочных данных |
| `internal/modules/product/assembler.go` | Сборка и преобразование данных товаров в структуры для API |
| `internal/modules/product/model.go` | Модели товаров, групп товаров и структур ответов |
| `internal/modules/product/repository.go` | Интерфейс получения отдельных товаров и связанных товаров по группе |
| `internal/modules/products/filter.go` | Построение SQL-условий фильтрации товаров с поддержкой исключения отдельных параметров при расчёте фасетов |
| `internal/modules/products/handler.go` | HTTP-обработчики запросов каталога, фильтрации и метаданных |
| `internal/modules/products/helpers.go` | Разбор и валидация параметров фильтрации из URL, обработка характеристик и вспомогательные функции для работы с фасетами и SQL-параметрами |
| `internal/modules/products/model.go` | Структуры параметров фильтрации, фасетов и метаданных фильтров |
| `internal/modules/products/repository.go` | Интерфейс и реализация запросов к каталогу, фильтрам и фасетам |
| `internal/modules/products/service.go` | Бизнес-логика получения товаров, применения фильтров и формирования метаданных |