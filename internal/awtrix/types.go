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
	// 0 = Icon doesn't move. 1 = Icon moves with text and will not appear again. 2 = Icon moves with text but appears again when the text starts to scroll again.
	PushIcon int `json:"pushIcon,omitempty"`
	// Sets how many times the text should be scrolled through the matrix before the app ends
	Repeat int `json:"repeat"`
	// Sets how long the app or notification should be displayed
	Duration int `json:"duration"`
	// Enables or disables autoscaling for bar and linechart
	Autoscale bool `json:"autoscale"`
	//  Defines the position of your custompage in the loop, starting at 0 for the first position. This will only apply with your first push. This function is experimental
	Pos int `json:"pos,omitempty"`
	// Removes the custom app when there is no update after the given time in seconds
	Lifetime int `json:"lifetime,omitempty"`
	// Defines if the **notification** will be stacked. false will immediately replace the current notification
	Stack bool `json:"stack"`
	// If the Matrix is off, the notification will wake it up for the time of the notification
	Wakeup bool `json:"wakeup"`
	// Disables the textscrolling
	NoScroll bool `json:"noScroll,omitempty"`
	// Modifies the scrollspeed. You need to enter a percentage value
	ScrollSpeed int `json:"scrollSpeed,omitempty"`
	// Shows an (https://blueforcer.github.io/awtrix-light/#/effects) as background
	Effect string `json:"effect,omitempty"`
	// Changes color and speed of the (https://blueforcer.github.io/awtrix-light/#/effects)
	EffectSettings map[string]any `json:"effectSettings,omitempty"`
}
