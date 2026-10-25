export function isComposingEvent(event) {
    return Boolean(event.isComposing || event.nativeEvent?.isComposing || event.keyCode === 229);
}
