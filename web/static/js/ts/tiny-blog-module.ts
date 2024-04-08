// import * as MD from './drawdown/drawdown.js'

class MarkDown extends AComponent{

    private getMarkDown(source: string) {
        return fetch(source)
            .then(response => response.text())
            .then((responseText) => {
                return responseText;
            });
    }

    public override initiate(element: HTMLElement) {
        const source = element.attributes.getNamedItem("source").value;
        this.getMarkDown(source)
            .then(val => {
                const div = document.createElement("div");
                element.insertAdjacentElement("afterend", div);
                div.innerHTML = markdown(val);
            });
        this.registerForCleanUp(element);
    }

    public override render() {
        return "no-op";
    }
}
