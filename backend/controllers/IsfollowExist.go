package controllers

import (
	"net/http"

	"rtf/help"
)

func (c *Controller) IsfollowExist(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		help.RespondNotOK(w, "notallowed")
		return
	}

	sender, ok := r.Context().Value("userID").(string)
	if !ok || sender == "" {
		help.RespondNotOK(w, "unauthorized")
		return
	}

	q := r.URL.Query()
	receiver := q.Get("receiver")

	isfollow := c.DB.IsFollowExist(sender, receiver)

	res := map[string]any{
		"isfollow": isfollow,
	}

	help.RespondOK(w, res, "followcheck")
}
