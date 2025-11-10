# Dusted Codes Blog & Website

Visit: https://dusted.codes

My personal website, blog and home of Dusted Codes Limited (UK based software development business).

Migrated from my previous [F# blog](https://github.com/dustinmoris/dustedcodes).

# About

My personal take on a simple Go blog engine which renders static HTML pages.

Blog articles can be written in Markdown or HTML. Markdown pages get compiled into static HTML during startup.

Feel free to fork it and create your own nerdy space in the world wide web!

# Cloudflare hosted CDN

I use Cloudflare R2 storage buckets and their CDN feature to host static assets behind https://cdn.dusted.codes.

Files inside `./cdn` mirror what goes into the CDN.

Upload contents into R2:

```bash
rclone copy ./cdn cf-dusted-codes:dusted-codes-cdn
```

Delete uploaded content (e.g. an entire folder):

```bash
rclone delete cf-dusted-codes:dusted-codes-cdn/folder-to-delete
```