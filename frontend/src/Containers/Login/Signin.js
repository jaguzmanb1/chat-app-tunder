import React, { Component } from 'react';
import './Signin.css';
import { Redirect } from 'react-router-dom';

export default class Signin extends Component {
    constructor(props){
        super(props);
        this.state = {
            phoneNumber: '',
            password: '',
            token: {},
            redirect: "false"
        }
        this.handlePhoneChange = this.handlePhoneChange.bind(this);
        this.handlePasswordChange = this.handlePasswordChange.bind(this);
        this.signin = this.signin.bind(this);
    }

    componentDidMount() {
        this.setState({
            redirect: localStorage.getItem('auth')

        })
    }

    handlePhoneChange(event) {
        this.setState({
            phoneNumber: event.target.value
        })
    }

    handlePasswordChange(event) {
        this.setState({
            password: event.target.value
        })
    }

    signin () {
        fetch('https://localhost:3031/signin', {
            method: 'POST',
            body: JSON.stringify({
                phone: this.state.phoneNumber,
                password: this.state.password,
            })
        })
        .then(async response =>{
            const data = await response.json()

            localStorage.setItem('token', data.message);
            localStorage.setItem('phone', this.state.phoneNumber);
            localStorage.setItem('auth', "true");

            this.setState( {
                redirect: "true"
            })
            window.location.reload(false);

            //console.log(localStorage.getItem('token'))
        })
        .catch((error) => {
            console.error(error);
        })      
    }

    render () {
        return (
            (
                this.state.redirect == "true" ? 
                    <Redirect to="/chat"/>
                : 
                    <div>
                        <div class="main">
                            <p class="sign" align="center">Sign in</p>
                            <div class="form1">
                                <input class="un" type="text" align="center" placeholder="Username" value={this.state.phoneNumber} onChange={this.handlePhoneChange}/>
                                <input class="pass" type="password" align="center" placeholder="Password" value={this.state.password} onChange={this.handlePasswordChange}/>
                                <button class="submit" align="center" onClick={this.signin}>Sign in</button>
                                <p class="forgot" align="center"><a href="#">Forgot Password?</a></p> 
                            </div>   
                        </div>
                    </div>
                )

            );
    }
}
