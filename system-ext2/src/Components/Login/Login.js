import Container from 'react-bootstrap/Container';
import Form from 'react-bootstrap/Form';
import Button from 'react-bootstrap/Button';
import './login.css'
const Login = () => {
    return (
        <Container id="main-container" className="d-grid h-100">
            <Form id="sign-in-form" className="text-center p-3 w-100">
                <h1 className="mb-3 fs-3 fw-normal">Login</h1>
                <Form.Group controlId="id-partition">
                    <Form.Control type="id" size="lg" placeholder="ID Partición" autoComplete="id" className="position-relative" />
                </Form.Group>
                <Form.Group controlId="user">
                    <Form.Control type="user" size="lg" placeholder="User" autoComplete="username" className="position-relative" />
                </Form.Group>
                <Form.Group className="mb-3" controlId="password">
                    <Form.Control type="password" size="lg" placeholder="Password" autoComplete="current-password" className="position-relative" />
                </Form.Group>
                <div className="d-grid">
                    <Button variant="primary" size="lg">Sign in</Button>
                </div>
                <p className="mt-5 text-muted">&copy; File System - EXT2</p>
            </Form>
        </Container>
    )
}

export default Login