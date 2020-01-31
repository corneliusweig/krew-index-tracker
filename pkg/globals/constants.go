package globals

const (
	// Your Google Cloud Platform project ID.
	ProjectID = "krew-k8s"

	// The name of the BigQuery dataset that contains the tables
	BQDataset = "plugin_download_count"

	// The default number of retries for scraping and uploading
	DefaultRetries = 5
)
