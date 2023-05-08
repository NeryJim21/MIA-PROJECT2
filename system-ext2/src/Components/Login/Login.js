import Container from 'react-bootstrap/Container';
import Form from 'react-bootstrap/Form';
import Button from 'react-bootstrap/Button';
import './login.css'
const Login = () => {

    function login(){
        var texto = "login >pwd=" + document.getElementById('inputPassword').value + " >user=" + document.getElementById('inputUser').value + " >id=" + document.getElementById('inputIdParticion').value
        var obj = { 'cmd':  texto}
        fetch(`http://localhost:5000/ejecutar`, {
          method: 'POST',
          body: JSON.stringify(obj),
        })
          .then(res => res.json())
          .catch(err => {
            console.error('Error:', err)
            alert("Ocurrio un error, ver la consola")
          })
          .then(response => {
            alert("Respuesta: " + response.result)
          })
    }


    return (
        <Container id="main-container" className="d-grid h-100">
            <Form id="sign-in-form" className="text-center p-3 w-100">
                <h1 className="mb-3 fs-3 fw-normal">Login</h1>
                <Form.Group controlId="id-partition">
                    <Form.Control type="id" id="inputIdParticion" size="lg" placeholder="ID Partición" autoComplete="id" className="position-relative" />
                </Form.Group>
                <Form.Group controlId="user">
                    <Form.Control type="user" id="inputUser" size="lg" placeholder="User" autoComplete="username" className="position-relative" />
                </Form.Group>
                <Form.Group className="mb-3" controlId="password">
                    <Form.Control type="password" id="inputPassword" size="lg" placeholder="Password" autoComplete="current-password" className="position-relative" />
                </Form.Group>
                <div className="d-grid">
                    <Button variant="primary" size="lg"onClick={login}>Sign in</Button>
                </div>
                <p className="mt-5 text-muted">&copy; File System - EXT2</p>
            </Form>
        </Container>
    )
}

export default Login