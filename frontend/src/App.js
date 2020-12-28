import Routes from './Routes';
import { BrowserRouter as Router } from 'react-router-dom';
import NavBar from './Components/NavBar'
import 'bootstrap/dist/css/bootstrap.min.css';

function App() {

  return (
    <div className="App"> 
      <NavBar/>
      <Router>
        <Routes/>
      </Router>
    </div>
  );
}

export default App;
