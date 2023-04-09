package messages

import "gtihub.com/televi-go/televi/telegram"

type ProgressingRequest[TSuccess telegram.Request, TFallback telegram.Request] struct {
	Success    TSuccess
	Fallback   TFallback
	IsFallback bool
}

func MakeProgressive[TRequest telegram.Request](request TRequest) ProgressingRequest[TRequest, TRequest] {
	return ProgressingRequest[TRequest, TRequest]{
		Success:    request,
		Fallback:   request,
		IsFallback: false,
	}
}

func Progress[
	TFurtherSuccess telegram.Request,
	TFurtherFallback telegram.Request,
	TSuccess telegram.Request,
	TFallback telegram.Request,
](
	request ProgressingRequest[TSuccess, TFallback],
	successMap func(TSuccess) TFurtherSuccess,
	fallbackMap func(TFallback) TFurtherFallback,
) ProgressingRequest[TFurtherSuccess, TFurtherFallback] {
	if request.IsFallback {
		return ProgressingRequest[TFurtherSuccess, TFurtherFallback]{
			Fallback:   fallbackMap(request.Fallback),
			IsFallback: true,
		}
	}
	return ProgressingRequest[TFurtherSuccess, TFurtherFallback]{
		Success:    successMap(request.Success),
		IsFallback: false,
	}
}
