package uploads

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func UploadFileHandler(c echo.Context, storage *BlobStore) error {
	var res ResourceRequest
	if err := c.Bind(&res); err != nil {
		return present(c, http.StatusBadRequest, nil, err)
	}

	err := storage.StoreRawImage(res.Name, "0", res.File)
	if err != nil {
		return present(c, http.StatusBadRequest, nil, err)
	}

	return c.String(http.StatusCreated, "File successfully uploaded!")
}

func ServeFileHandler(c echo.Context, storage *BlobStore) error {
	req := c.Request()
	requestedFile := strings.TrimPrefix(req.URL.Path, "/serve/")

	image, err := storage.GetImage(requestedFile, hasSize(req))
	if err != nil {
		return err
	}
	return c.Blob(http.StatusOK, "image/jpeg", image)
}

func hasSize(req *http.Request) string {
	params := req.URL.Query()
	return params.Get("size")
}

func present(c echo.Context, status int, pl interface{}, err error) error {
	if err == nil {
		return c.JSON(status, pl)
	} else {
		return c.JSON(status, echo.Map{"error": err.Error()})
	}
}
