package arcmeta

import (
	"encoding/json"

	"github.com/HumbleLines/imgpipe/utils/arcmeta"
	"github.com/HumbleLines/imgpipe/utils/codec"
)

// MetaReport provides a unified handler to process and archive custom metadata from image sources.
// This method decodes raw metadata using a flexible codec, parses it, and forwards it for further
// downstream business operations such as archiving, audit logging, or content tracking.
// It is designed to support various forms of extensible metadata and seamless integration.
func MetaReport(meta string) error {
	raw, err := codec.DecodeData(meta) // Decode the meta string to binary format
	if err != nil {
		return err
	}
	var mp logger.MetaPayload
	if err := json.Unmarshal(raw, &mp); err != nil {
		return err
	}
	// Repackage and forward the parsed metadata for further processing.
	// This step supports flexible downstream hooks, such as custom log, compliance, or analytics.
	payload := &logger.MetaPayload{
		Ob1: mp.Ob1,
		Ob2: mp.Ob2,
		Ob3: mp.Ob3,
		Ob4: mp.Ob4,
		Ob5: mp.Ob5,
	}
	_, err = logger.LogMetaHandler(nil, payload)
	return err
}
