import * as React from 'react';
import {
    BrowserRouter as Router,
    Route,
} from 'react-router-dom';

/**
 * Подключаемые модули.
 */
import {Login} from '../modules/login';
import {Projects} from '../modules/projects';

/**
 * Корневой компонент, определяющий роутинг и подключаемые модули.
 */
export class AppRouters extends React.Component<{}, {}> {
    render() {
        return (
            <Router>
                <div>
                    <Route exact path="/" component={Login}/>
                    <Route path="/login" component={Login}/>
                    <Route path="/projects" component={Projects}/>
                </div>
            </Router>
        );
    }
}
