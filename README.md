# IoT Final Project: Высоконагруженный сервис с аналитикой на Go в Kubernetes

## Описание проекта

Этот проект реализует высоконагруженный сервис на языке Go для обработки потоковых метрик от IoT-устройств (симуляция данных: timestamp, CPU, RPS). Сервис предназначен для мониторинга нагрузки в реальном времени, с простой статистической аналитикой (rolling average для сглаживания и z-score для детекции аномалий).

**Реализованно:**
- **Core-сервис на Go:** HTTP-эндпоинты (/metrics POST для приема JSON-метрик, /analyze GET для аналитики, /health GET). Обработка >1000 RPS. Goroutines и channels для асинхронной аналитики.
- **Кэширование:** Redis - метрики сохраняются с ключом "metric:timestamp" и TTL 5 мин.
- **Аналитика:** Rolling average (окно 50 событий) для сглаживания, z-score детекция аномалий (threshold=2.0σ). Mutex для thread-safety.
- **Мониторинг:** Prometheus - counters, histogram. Экспорт на /metrics/prometheus.
- **Развертывание в Kubernetes:** Minikube/Kind или облако. Deployment (2 реплики + Redis), Service, HPA, Ingress. Prometheus ConfigMap.
- **Тестирование:** Locust для нагрузки, Grafana dashboard.
- **Оптимизация:** Multi-stage Dockerfile, ресурсы в K8s.



**Стек технологий:**
- Go 1.24 (gorilla/mux, go-redis/v8, prometheus/client_golang).
- Redis 7-alpine.
- Kubernetes (apps/v1, autoscaling/v2, networking.k8s.io/v1).
- Prometheus + Grafana (базовые дашборды).
- Locust для нагрузки.

## Мониторинг

- Prometheus: Доступ на http://<minikube-ip>:9090 .
- Grafana: Дашборд с панелями для RPS, anomalies, latency (импорт /k8s/grafana-dashboard.json).
- Аномалии: Если >5/мин - manual alert.

## Архитектура

- **Поток данных:** IoT → POST /metrics (JSON) → Redis cache → Analyzer → Prometheus metrics.
- **K8s:** Pods (Go + Redis), HPA scales на CPU, Ingress routes.

## Оптимизация и производительность

- Latency <50ms под 1000 RPS (тесты).
- Redis: TTL 5min, no eviction tune.
- Go: pprof для профилирования (curl http://localhost:8080/debug/pprof/).
- Scaling: HPA растет до 5 pods при >70% CPU.
