package providers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mehmetymw/search-aggregation-service/backend/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestXmlProviderClient_FetchContents(t *testing.T) {
	mockResponse := `
	<feed>
		<items>
			<item>
				<id>2</id>
				<headline>Test Article</headline>
				<type>article</type>
				<stats>
					<views>50</views>
					<likes>5</likes>
					<duration></duration>
					<reading_time>300</reading_time>
					<reactions>20</reactions>
					<comments>1</comments>
				</stats>
				<publication_date>2023-10-26</publication_date>
				<categories>
					<category>news</category>
					<category>tech</category>
				</categories>
			</item>
		</items>
	</feed>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	client := NewXmlProviderClient()
	provider := entity.Provider{
		BaseURL: server.URL,
	}

	ctx := context.Background()
	items, err := client.FetchContents(ctx, provider)

	assert.NoError(t, err)
	assert.Len(t, items, 1)

	item := items[0]
	assert.Equal(t, "2", item.ProviderContentID)
	assert.Equal(t, "Test Article", item.Title)
	assert.Equal(t, "article", item.ContentType)
	assert.Equal(t, int32(300), item.ReadingTime)
	assert.Len(t, item.Tags, 2)
}

func TestXmlProviderClient_FetchContents_DurationParsing(t *testing.T) {
	mockResponse := `
	<feed>
		<items>
			<item>
				<id>3</id>
				<headline>Video</headline>
				<type>video</type>
				<stats>
					<duration>2:30</duration>
				</stats>
				<publication_date>2023-10-26</publication_date>
			</item>
		</items>
	</feed>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	client := NewXmlProviderClient()
	provider := entity.Provider{BaseURL: server.URL}

	items, err := client.FetchContents(context.Background(), provider)
	assert.NoError(t, err)
	assert.Equal(t, int32(150), items[0].DurationSec) // 2*60 + 30 = 150
}
