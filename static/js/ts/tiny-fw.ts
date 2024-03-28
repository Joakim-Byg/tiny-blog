
const mwfwLib = new Map<string, IComponent>();

function initTags() {
    mwfwLib.set("my-tag", new MyTag());
}
function renderTags() {
    for (let key of Array.from(mwfwLib.keys())) {
        console.log(`went here: ${key}`);
        let elements = Array.from(document.getElementsByTagName(key));
        console.log(`Element size: ${elements.length}`);
        for (let elem of elements) {
            let _tag = mwfwLib.get(key);
            _tag.initiate(elem as HTMLElement);
        }
    }
}

interface IComponent {
    initiate(element: HTMLElement):void;
    render():string;
}

document.onreadystatechange = function () {
    if (document.readyState === 'complete') {
        initTags();
        renderTags();
    }
}

