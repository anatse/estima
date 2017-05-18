import * as React from 'react';

/**
 * TODO Здесь будет главная страница с проектами.
 */
export class Projects extends React.Component<any, void> {

    componentWillMount () {
        document.title = 'Проекты';
    }

    render() {
        return (
            <div>
                Страница проектов.
            </div>
        );
    }
}
