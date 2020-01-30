import React from 'react'
import AudioRecorder from './components/AudioRecorder'

require('./App.scss')
require('./scss-lib/normalize.scss')
require('./scss-lib/font-awesome/css/all.min.css')
function App() {
  return (
    <div>
      <AudioRecorder/>
    </div>
  );
}

export default App;
