
const AUTH_HEADER = "authorization"

const YT_DESKTOP_API_URL = "http://localhost:9863/api/v1/"
const YT_DESKTOP_GET_STATE_URL = YT_DESKTOP_API_URL + "state"
const YT_VIDEO_URL = "https://www.youtube.com/watch?v="

const SECONDS_IN_HOUR = 3600
const SECONDS_IN_MINUTE = 60


function formatTime(seconds) {
    let formattedTime = ""

    let hours = parseInt(seconds / SECONDS_IN_HOUR)
    if (hours != 0) {
        formattedTime += `${hours}ч `
        seconds = seconds % SECONDS_IN_HOUR
    }

    let minutes = parseInt(seconds / SECONDS_IN_MINUTE)
    if (minutes != 0) {
        formattedTime += `${minutes}м `
        seconds = seconds % SECONDS_IN_MINUTE
    }

    if (seconds != 0) {
        formattedTime += `${seconds}с `
    }

    return formattedTime
}

function getCurrentSongInfo(video, player) {
    const videoUrl = YT_VIDEO_URL + video.Id
    const timestamp = formatTime(parseInt(player.videoProgress))

    return `${video.title} ${videoUrl} таймстемп: ${timestamp}`
}

function showResult(text) {
    document.querySelector("html").innerHTML = text
}


async function main() {
    const urlParams = new URLSearchParams(window.location.search);
    const authToken = urlParams.get("token")

    if (!authToken) {
        showResult("Не могу получить информацию о текущей песне без токена авторизации")
        return
    }

    try {
        const res = await fetch(YT_DESKTOP_GET_STATE_URL, {
            headers: {
                [AUTH_HEADER]: authToken
            }
        })

        const ytState = await res.json()

        if (ytState.error) {
            showResult(ytState.error)
            return
        }

        showResult(getCurrentSongInfo(ytState.video, ytState.player))
    } catch (err) {
        console.log(err)
        showResult("Что-то пошло не так во время получения информации о текущей песне")
    }
}

main()
