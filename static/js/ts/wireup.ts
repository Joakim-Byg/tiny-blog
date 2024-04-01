
const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms));

document.onreadystatechange = function () {
    if (document.readyState === 'complete') {
        registerTag("my-tag", new MyTag());
        registerTag("md", new MarkDown());
        renderTags();
        sleep(500).then(value => cleanupElements());
    }
}