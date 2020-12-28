import React, { Component } from 'react';
import { Redirect } from 'react-router-dom';

export default class Signout extends Component {
    constructor(props){
        super(props);


    }
    componentDidMount() {
        localStorage.setItem('token', "");
        localStorage.setItem('phone', "");
        localStorage.setItem('auth', "");
        this.props.history.push('/');
        window.location.reload(false);

    }

    render () {
        return (
            
                <div>
                    

                </div>
        )
    }
}
