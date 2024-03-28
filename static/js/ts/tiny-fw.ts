
const tfwTags = new Map<string, IComponent>();
const renderCleanups = new Array<HTMLElement>();

function initTags() {
    tfwTags.set("my-tag", new MyTag());
}
function renderTags() {
    for (let key of Array.from(tfwTags.keys())) {
        console.log(`went here: ${key}`);
        let elements = Array.from(document.getElementsByTagName(key));
        console.log(`Element size: ${elements.length}`);
        for (let elem of elements) {
            let _tag = tfwTags.get(key);
            _tag.initiate(elem as HTMLElement);
        }
    }
}

function cleanupElements() {
    for (let elem of renderCleanups) {
        elem.remove();
    }
}

interface IComponent {
    initiate(element: HTMLElement):void;
    render():string;
    registerForCleanUp(element: HTMLElement):void;
}

abstract class AComponent implements IComponent {
    abstract initiate(element: HTMLElement): void;
    abstract render(): string;
    registerForCleanUp(element: HTMLElement): void {
        renderCleanups.push(element);
    }

}

document.onreadystatechange = function () {
    if (document.readyState === 'complete') {
        initTags();
        renderTags();
        cleanupElements();
    }
}

