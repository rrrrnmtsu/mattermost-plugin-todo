import {configureLocale, translate} from './localization';

let locale = 'en';
const currentLocale = () => locale;

describe('current profile language', () => {
    test('follows the current profile instead of the browser language', () => {
        locale = 'ja-JP';
        configureLocale(currentLocale);
        expect(translate('Tasks', 'タスク')).toBe('タスク');
        locale = 'en';
        expect(translate('Tasks', 'タスク')).toBe('Tasks');
    });
});
