const message = (name, ...detail) => {
    const parsedDetail = detail.length ? detail.join(',') : ''
    const payload = `${btoa(name)},${btoa(parsedDetail)}`

    return new Promise(res => {
        window.__resolveMessage = response => res(JSON.parse(atob(response)))
        window.webkit.messageHandlers.firecrypt.postMessage(payload)
    })
}

if (window.__FIRECRYPT_RELEASE) {
    window.addEventListener('contextmenu', event => {
        event.preventDefault()
    })
} else {
    console.info('this is a debug build. issues not present in release builds may occur.')
}
