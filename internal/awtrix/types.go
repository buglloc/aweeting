package awtrix

type Payload struct {
	// The text to display
	Text string `json:"text"`
	// Changes the Uppercase setting. 0=global setting, 1=forces uppercase; 2=shows as it sent
	TextCase int `json:"textCase,omitempty"`
	// Draw the text on top
	TopText bool `json:"topText,omitempty"`
	//Sets an offset for the x position of a starting text
	TextOffset int `json:"textOffset,omitempty"`
	// The text, bar or line color (#hex)
	Color string `json:"color,omitempty"`
	// Sets a background color (#hex)
	Background string `json:"background,omitempty"`
	// Fades each letter in the text differently through the entire RGB spectrum
	Rainbow bool `json:"rainbow,omitempty"`
	// The icon ID or filename (without extension) to display on the app
	Icon string `json:"icon,omitempty"`
}
