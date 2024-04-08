
const tfwTags = new Map<string, IComponent>();
const renderCleanups = new Array<HTMLElement>();


function renderTags() {
    for (let key of Array.from(tfwTags.keys())) {
        let elements = Array.from(document.getElementsByTagName(key));
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

function registerTag(tag: string, component:AComponent) {
    tfwTags.set(tag, component);
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

