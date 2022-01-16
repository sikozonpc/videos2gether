import React from 'react'
import { Router, Route, Switch } from 'react-router'
import { createBrowserHistory } from 'history'
import Room from './screens/room/Room'
import Hub from './screens/hub/Hub'

const history = createBrowserHistory()

const AppRouter: React.FC = () => {
  return (
    <Router history={history}>
      <Switch>
        <Route path='/' component={Hub} exact />
        <Route path='/room/:roomID' component={Room} exact />
        <Route component={() => <h1>Not found</h1>} />
      </Switch>
    </Router>
  )
}

export default AppRouter