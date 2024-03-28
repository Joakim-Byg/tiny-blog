
class MyTag implements IComponent{
    message:string;

    public initiate(element: HTMLElement) {
        this.message = element.attributes.getNamedItem("message").value;
        element.insertAdjacentHTML("afterend", this.render());
        element.remove();
    }

    public render() {
        return `<p>${this.message}</p>`;
    }
}
