# Youtube2NetDisk

+ A Youtube - Baidu Net Disk Auto Carrier.

⚠️ Downloading video from Youtube could violate [Youtube's term of service](https://www.youtube.com/t/terms),

> Content is provided to you AS IS. You may access Content for your information and personal use solely as intended through the provided functionality of the Service and as permitted under these Terms of Service. You shall not download any Content unless you see a “download” or similar link displayed by YouTube on the Service for that Content. You shall not copy, reproduce, distribute, transmit, broadcast, display, sell, license, or otherwise exploit any Content for any other purposes without the prior written consent of YouTube or the respective licensors of the Content. YouTube and its licensors reserve all rights not expressly granted in and to the Service and the Content.

**Use at your own risk!**
**You should make sure you have granted permission from the channel owner before you start carrying videos.**

## Thanks to

This program is powered by these awesome projects, thanks to them and their wonderful author!

+ [youtube-dl](https://github.com/rg3/youtube-dl): Command-line program to download videos from YouTube.com and other video sites
+ [bypy](https://github.com/houtianze/bypy): Python client for Baidu Yun (Personal Cloud Storage)

## Legacy Version

For the legacy version, checkout branch [legacy](https://github.com/cool2645/youtube-to-netdisk/tree/legacy).

## How to use

```go
package main

import (
	"github.com/cool2645/youtube-to-netdisk/carrier"
	"github.com/cool2645/youtube-to-netdisk/http"
)

func main() {
	carrier.Use(http.WebInterface{})
	carrier.Start()
}

```

In the code above you start a carrier with a web interface.  
You'll be able to trigger and cancel tasks via a web interface.

## http.WebInterface

### API Reference

Refer to the APIs in case you need.

**POST: /api/trigger**
+ Content-Type: application/x-www-form-urlencoded
+ Description: Trigger a new carrier task

    | Parameter | Description |
    | --- | --- |
    | url | The youtube video url. |
    | title | The youtube video title. |
    | author_name | The author's name of the youtube video. |
    | description | The description of the youtube video. |

**DELETE: /api/cancel/{id}**
+ Content-Type: application/x-www-form-urlencoded
+ Description: Cancel a task.
+ Args: {id}: ID of the task.

**GET/POST: /api/retry/{id}**
+ Content-Type: application/x-www-form-urlencoded
+ Description: Retry a task.
+ Args: {id}: ID of the task.

**GET: /api/task**
+ Content-Type: application/x-www-form-urlencoded
+ Description: Get tasks (Log is not shown).

    | Parameter | Description |
    | --- | --- |
    | state | "Rejected" if you want rejected tasks, otherwise rejected tasks are excluded. |
    | order | "desc" or "asc" |
    | page | Page num. |
    | perPage | Count of entries in each page. |

**GET: /api/task/{id}**
+ Content-Type: application/x-www-form-urlencoded
+ Description: Get a specific task.
+ Args: {id}: ID of the task.

**GET: /api/task/{id}/log**
+ Content-Type: application/x-www-form-urlencoded
+ Description: Get log of a specific task.
+ Args: {id}: ID of the task.

**GET: /api/keyword**
+ Content-Type: application/x-www-form-urlencoded
+ Description: Show all keywords.

**POST: /api/keyword**
+ Content-Type: application/x-www-form-urlencoded
+ Description: Add a keyword.

    | Parameter | Description |
    | --- | --- |
    | keyword | Keyword |

**DELETE: /api/keyword**
+ Content-Type: application/x-www-form-urlencoded
+ Description: Delete one keyword.

    | Parameter | Description |
    | --- | --- |
    | keyword | Keyword |

### Build Front-End

Build Front-end so that you can see tasks and running processes from the website.

cd to view/, and then run

```
yarn install
yarn build
```

### IFTTT

Now that we can trigger one task to download and upload video for us through the web interface (refer to API doc above).

With [IFTTT](https://ifttt.com)'s awesome service we can watch a subscribed channel and trigger tasks automatically (use the Youtube and WebHook service) when a new video is published.

However in my scenario I do not always need the videos, so here the "keywords" are. Only videos with keywords in the title will be carried.
For instance, according to my case, the keywords are "ラブライブ" and "虹ヶ咲", and one of my subscribed channels is "電撃オンライン", which is likely to also publish videos that has nothing to do with the topic I want.

If you don't need this "keyword" feature, emmmm, maybe you can just add a blank keyword?

## Contribute

Feel free to contribute by Issue and Pull Requests!

Made with ♥️ by 梨子. Theme by Bittersweet. Thanks to you.
