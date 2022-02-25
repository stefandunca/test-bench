package main

import (
	"fmt"

	"github.com/linkedin/goavro/v2"
)

func main() {
	codec, err := goavro.NewCodec(`
		{
			"type": "record",
			"name": "LongList",
			"fields" : [
				{
					"name": "next",
				 	"type": [
					 	"null", "LongList",
						{
							"type": "long",
							"logicalType": "timestamp-millis"
						}
				],
					 "default": null
				}
			]
		}`)
	if err != nil {
		fmt.Println(err)
	}
	textual := []byte(`{"next":{"LongList":{}}}`)
	native, _, err := codec.NativeFromTextual(textual)
	if err != nil {
		fmt.Println(err)
	}
	// Convert native Go form to binary Avro data
	binary, err := codec.BinaryFromNative(nil, native)
	if err != nil {
		fmt.Println(err)
	}

	// Convert binary Avro data back to native Go form
	native, _, err = codec.NativeFromBinary(binary)
	if err != nil {
		fmt.Println(err)
	}

	// Convert native Go form to textual Avro data
	textual, err = codec.TextualFromNative(nil, native)
	if err != nil {
		fmt.Println(err)
	}

	// NOTE: Textual encoding will show all fields, even those with values that
	// match their default values
	fmt.Println(string(textual))
	// Output: {"next":{"LongList":{"next":null}}}
}
