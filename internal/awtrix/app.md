Source: https://raw.githubusercontent.com/Blueforcer/awtrix-light/main/docs/api.md

| Key | Type | Description | Default | Custom App | Notification |
| --- | ---- | ----------- | ------- | ------- |------- |
| `text` | string | The text to display. | N/A | X | X |
| `textCase` | integer | Changes the Uppercase setting. 0=global setting, 1=forces uppercase; 2=shows as it sent. | 0 | X | X |
| `topText` | boolean | Draw the text on top | false | X | X |
| `textOffset` | integer | Sets an offset for the x position of a starting text. | 0 | X | X |
| `color` | string or array of integers | The text, bar or line color | N/A | X | X |
| `background` | string or array of integers | Sets a background color | N/A | X | X |
| `rainbow` | boolean | Fades each letter in the text differently through the entire RGB spectrum. | false | X | X |
| `icon` | string | The icon ID or filename (without extension) to display on the app. | N/A | X | X |
| `pushIcon` | integer | 0 = Icon doesn't move. 1 = Icon moves with text and will not appear again. 2 = Icon moves with text but appears again when the text starts to scroll again. | 0 | X | X |
| `repeat` | integer | Sets how many times the text should be scrolled through the matrix before the app ends. | 1 | X | X |
| `duration` | integer | Sets how long the app or notification should be displayed. | 5 | X | X |
| `hold` | boolean | Set it to true, to hold your **notification** on top until you press the middle button or dismiss it via HomeAssistant. This key only belongs to notification. | false |   | X |
| `sound` | string | The filename of your RTTTL ringtone file placed in the MELODIES folder (without extension). | N/A |   | X |
| `rtttl` | string | Allows to send the RTTTL sound string with the json | N/A |   | X |
| `loopSound` | boolean | Loops the sound or rtttl as long as the notification is running | false |   | X |
| `bar` | array of integers | draws a bargraph. Without icon maximum 16 values, with icon 11 values | N/A | X | X |
| `line` | array of integers | draws a linechart. Without icon maximum 16 values, with icon 11 values | N/A | X | X |
| `autoscale` | boolean | Enables or disables autoscaling for bar and linechart | true | X | X |
| `progress` | integer | Shows a progressbar. Value can be 0-100 | -1 | X | X |
| `progressC` | string or array of integers  | The color of the progressbar | -1 | X | X |
| `progressBC` | string or array of integers  | The color of the progressbar background | -1 | X | X |
| `pos` | integer | Defines the position of your custompage in the loop, starting at 0 for the first position. This will only apply with your first push. This function is experimental | N/A |  X |   | 
| `draw` | array of objects | Array of drawing instructions. Each object represents a drawing command. | See the drawing instructions below | X | X |
| `lifetime` | integer | Removes the custom app when there is no update after the given time in seconds | 0 | X |   |
| `stack` | boolean | Defines if the **notification** will be stacked. false will immediately replace the current notification | true |   | X |
| `wakeup` | boolean | If the Matrix is off, the notification will wake it up for the time of the notification. | false |   | X |
| `noScroll` | boolean | Disables the textscrolling | false | X | X |
| `clients` | array of strings | Allows to forward a notification to other awtrix. Use the MQTT prefix for MQTT and IP adresses for HTTP |  |   | X |
| `scrollSpeed` | integer | Modifies the scrollspeed. You need to enter a percentage value | 100 | X | X |
| `effect` | string | Shows an [effect](https://blueforcer.github.io/awtrix-light/#/effects) as background |  | X | X |  
| `effectSettings` | json map | Changes color and speed of the [effect](https://blueforcer.github.io/awtrix-light/#/effects) |  | X | X |  
| `save` | boolean | Saves your customapp into flash and reload it after boot. You should avoid that with customapps wich has high update frequency because ESPs flashmemory has limited writecycles  |  | X |  |  