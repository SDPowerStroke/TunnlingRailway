const express = require("express");
const { spawn } = require("child_process");

const app = express();
const PORT = process.env.PORT || 3000;

// ⚠️ Hardcoded for demo only
const YT_URL = "https://www.youtube.com/watch?v=9kRhE5vgCvY";

app.get("/video", (req, res) => {
    res.setHeader("Content-Type", "video/mp4");

    const yt = spawn("python3", [
        "-m",
        "yt_dlp",
        "-f",
        "best[ext=mp4]/best",
        "-o",
        "-",
        YT_URL
    ]);

    yt.stdout.pipe(res);

    yt.stderr.on("data", d => console.error(d.toString()));

    yt.on("error", err => {
        console.error(err);
        res.status(500).end("yt-dlp error");
    });
});

app.use(express.static("public"));

app.listen(PORT, () => {
    console.log(`Listening on port ${PORT}`);
});
