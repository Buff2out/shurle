package logging

import (
	"go.uber.org/zap"
)

func GetSugaredLogger() (*zap.SugaredLogger, *zap.Logger) {
	// добавляем предустановленный логер NewDevelopment
	logger, err := zap.NewDevelopment()
	if err != nil {
		// вызываем панику, если ошибка
		panic("cannot initialize zap")
	}

	// делаем логер SugaredLogger
	sugar := logger.Sugar()
	return sugar, logger
}

//func GetStringMsg(urI string, msg string) string {
//
//	return fmt.Sprintf("%d", 10)
//}
//
//// бесполезнейшая функция, но зато мне быстрее можно вспомнить как вычитать одно время из другого.
//func DeltaTime(t1 time.Time, t2 time.Time) time.Duration {
//	return t2.Sub(t1)
//}
