
class MyTag extends AComponent{
    message:string;

    public initiate(element: HTMLElement) {
        this.message = element.attributes.getNamedItem("message").value;
        element.insertAdjacentHTML("afterend", this.render());
        this.registerForCleanUp(element);
    }

    public render() {
        return `<p>${this.message}</p>`;
    }
}

registerTag("my-tag", new MyTag());