import {isComposingEvent} from './keyboard';

describe('IME confirmation', () => {
    test.each([{isComposing: true}, {nativeEvent: {isComposing: true}}, {keyCode: 229}])('recognizes active composition %p', (event) => {
        expect(isComposingEvent(event)).toBe(true);
    });
    test('allows the following normal Enter', () => {
        expect(isComposingEvent({key: 'Enter', keyCode: 13})).toBe(false);
    });
});
