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

## How it works

The "youtube-to-netdisk" provides a web interface via which you can trigger a task that starts a process downloading as well as uploading the video to Baidu net disk.

## How to use

Before you start working with this program, read the following introduction carefully.

### Web API

Refer to these APIs in case you need.

**POST: /api/trigger**
+ Content-Type: application/x-www-form-urlencoded
+ Description: Trigger a new carrier task

    | Parameter | Description |
    | --- | --- |
    | url | The youtube video url. |
    | title | The youtube video title. |
    | author_name | The author's name of the youtube video. |
    | description | The description of the youtube video. |

**DELETE: /api/kill/{id}**
+ Content-Type: application/x-www-form-urlencoded
+ Description: Kill a running task.
+ Args: {id}: ID of the task.

**GET/POST: /api/retry/{id}**
+ Content-Type: application/x-www-form-urlencoded
+ Description: Retry a task.
+ Args: {id}: ID of the task.

**GET: /api/running**
+ Content-Type: application/x-www-form-urlencoded
+ Description: Get all running tasks.

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

### Python Script

The two python scripts are the core of the program. You can manually run them.

For instance,

```
python3 download.py https://youtu.be/yourlinkhere
```

This will download the video to current directory.

```
python3 syncBaidu.py filenamedownloadedjustnow.mp4 examplefolder
```

This will upload the video to the folder 我的应用数据/bypy/examplefolder in your Baidu net disk.

### Baidu Net Disk

If it is the first time you use bypy, you should manually run bypy or the script `syncBaidu.py`(check above for usage).
Bypy will ask you for authorization, afterwards it will be able to access your net disk (with saved credential).

### Server

Some explanation of configs.

+ **python_cmd** The python command. Usually it is `python3`, for some OS, like ArchLinux, it could be `python`.
+ **temp_path** The path where temp file for logs are stored.
+ **netdisk_folder** To which videos are uploaded.
+ **netdisk_sharelink** && **netdisk_sharepass** You can create a share link of the folder (parent of the video), after that the program can generate a share link for each video uploaded, e.g. https://pan.baidu.com/s/xxxxxx?fid=1111111

The web interface will be listening at **port**.

Videos will be stored at folder `static` (Perhaps should be created manually?)

### RiRi Notification

This branch uses [ritorudemonriri](https://github.com/rikakomoe/ritorudemonriri) as notification driver.

Turn on the **riri_enable** flag in config will enable this feature.

Fill in **riri_key** with your channel key, and **riri_addr** with your ritorudemonriri address.

Listening commands:

* **/carrier_subscribe** Set up subscription for this chat.
* **/carrier_subscribe --condense** Set up condensed subscription for this chat.
* **/carrier_unsubscribe** Suspend subscription for this chat.

### qwqq.pw

The program use [qwqq.pw](https://qwqq.pw) as its url shorten service, since the video
name are usually quite long to share.

### Build Front-End

Build Front-end so that you can see tasks and running processes from the website.

cd to app/, and then run

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
