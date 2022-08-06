export interface Toaster {
    type: 'INFO' | 'ERROR' | 'WARNING';
    text: string;
}