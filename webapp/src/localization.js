// Sakura deployment UI labels; read only the current user's locale.
let getLocale = () => navigator.language || 'en';

export function configureLocale(source) {
    getLocale = source;
}

export function translate(english, japanese) {
    return (getLocale() || 'en').toLowerCase().startsWith('ja') ? japanese : english;
}
