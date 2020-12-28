import Navbar from 'react-bootstrap/Navbar'
import Nav from 'react-bootstrap/Nav'
import NavDropdown from 'react-bootstrap/NavDropdown'
import Form from 'react-bootstrap/Form'
import FormControl from 'react-bootstrap/FormControl'
import Button from 'react-bootstrap/Button'

export default ({}) => {
    return (
        <Navbar bg="light" expand="lg">
            <Navbar.Brand href="/">React-Bootstrap</Navbar.Brand>
            <Navbar.Toggle aria-controls="basic-navbar-nav" />
            <Navbar.Collapse id="basic-navbar-nav">
                <Nav className="mr-auto">
                <Nav.Link href="/">Home</Nav.Link>
                {(localStorage.getItem("auth") == "true" ? 
                    <Nav.Link href="/chat">Chat</Nav.Link>
                    :
                    null
                )}
                </Nav>
                <Nav>
                    {(localStorage.getItem("auth") == "true" ? 
                        <Nav.Link href="/signout">Signout</Nav.Link>
                    :
                        <Nav.Link href="/signin">Signin</Nav.Link>
                    )}
                </Nav>
                <Form inline>
                <FormControl type="text" placeholder="Search" className="mr-sm-2" />
                <Button variant="outline-success">Search</Button>
                </Form>
            </Navbar.Collapse>
        </Navbar>
    )
}