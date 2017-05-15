import * as React from 'react';
import {
    HashRouter as Router,
    Route,
} from 'react-router-dom';

/**
 * Подключаемые модули.
 */
import {Authentication} from '../modules/Authentication';
import {Projects} from '../modules/projects';

/**
 * Корневой компонент, определяющий роутинг и подключаемые модули.
 */
export class AppRouters extends React.Component<any, void> {
    render() {
        return (
            <Router>
                <div>
                    <Route exact path="/" component={Authentication} />
                    <Route path="/authentication" component={Authentication} />
                    <Route path="/projects" component={Projects} />
                </div>
            </Router>
        );
    }
}
