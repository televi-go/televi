package telegram

type GetUpdatesRequest struct {
	Offset         int
	Limit          int
	Timeout        int
	AllowedUpdates []string
}

func (getUpdatesRequest GetUpdatesRequest) Method() string {
	return "getUpdates"
}

func (getUpdatesRequest GetUpdatesRequest) Params() (Params, error) {
	params := make(Params)
	params.WriteNonZero("offset", getUpdatesRequest.Offset)
	params.WriteNonZero("limit", getUpdatesRequest.Limit)
	params.WriteNonZero("timeout", getUpdatesRequest.Timeout)
	if len(getUpdatesRequest.AllowedUpdates) != 0 {
		params.WriteJson("allowed_updates", getUpdatesRequest.AllowedUpdates)
	}

	return params, nil
}
