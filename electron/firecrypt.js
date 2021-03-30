const astilectronReady = new Promise(res => {
    document.addEventListener('astilectron-ready', () => {
        res()
    })
})

const message = (msg, ...detail) => {
    const parsedDetail = detail.length > 0 ? detail.join(',') : ''
    const encodedMessage = `${btoa(msg)},${btoa(parsedDetail)}`
    return new Promise(res => {
        astilectronReady.then(() => {
            astilectron.sendMessage(encodedMessage, res)
        })
    })
}