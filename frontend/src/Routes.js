import { Route, Switch } from 'react-router-dom';
import AppliedRoute from './Components/AppliedRoute';
import PrivateRoute from './Components/PrivateRoute';
import Signin from './Containers/Login/Signin'
import Chat from './Containers/Chat/Chat'
import NotFound from './Containers/NotFound';
import Signout from './Containers/Signout/Signout'

export default () => (

    <Switch>
        <AppliedRoute path="/" exact component={NotFound}/>
        <AppliedRoute path="/signin" exact component={Signin}/>
        <PrivateRoute path="/chat" loggedIn={localStorage.getItem("auth")} component={Chat}/>
        <PrivateRoute path="/signout" loggedIn={localStorage.getItem("auth")} component={Signout}/>
        <Route component={NotFound} />
    </Switch>
)