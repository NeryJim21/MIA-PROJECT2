import { BrowserRouter as Router, Route, Routes } from 'react-router-dom'
import Navbar from './Components/Navbar/Navbar';
import Editor from './Components/Editor/Editor';
import Login from './Components/Login/Login';
import Reportes from './Components/Reportes/Reportes';


function App() {
	return (
		<div>
			<Router>
				<Navbar/>
				<Routes>
					<Route path='/' element={<Editor/>}/>
					<Route path='/Editor' element={<Editor/>}/>
					<Route path='/Login' element={<Login/>}/>
					<Route path='/Reportes' element={<Reportes/>}/>
				</Routes>
			</Router>
		</div>
	);
}

export default App;
