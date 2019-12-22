import os 
from flask import Flask, flash, request, redirect, url_for
from werkzeug.utils import secure_filename
from lib.musicAssembler import MusicAssembler

UPLOAD_FOLDER = "/voiceTempStorage"
ALLOWED_EXTENSION = {".mp3"}
CEPHFS = False # if CephFS is enabled or not. In development mode, we don't use it.
app = Flask(__name__)
app.config['UPLOAD_FOLDER'] = UPLOAD_FOLDER

def allowed_file(filename):
    return '.' in filename and \
        filename.rsplit('.',1)[1].lower() in ALLOWED_EXTENSION

def getRandomBeatFile(associatedVoicefilename):
    '''
    Connect to CephFS and retrieve a random beatFile.
    Then, get the binaries of the file and save it inside /beatTempStorage folder with the name
    `beat_<filenameUUID>.mp3`. Finally, return `beat_<filenameUUID>.mp3` as a string.
    '''
    uuid = associatedVoicefilename.split("_")[1].split(".")[0]
    if (CEPHFS):
        
    else:
        



@app.route('/upload/<uuid>',methods=['POST'])
def upload_file(uuid):
    if request.method == 'POST':
        # check if the post request has the file part
        if 'file' not in request.files:
            flash('No file part')
            return redirect(request.url)
        file = request.files['file']
        if file.filename == '':
            flash('No selected file')
            return redirect(request.url)
        if file and allowed_file(file.filename):
            filename = secure_filename("voice_"+uuid+".mp3") # replace the original name by its uuid
            file.save(os.path.join(app.config['UPLOAD_FOLDER'], filename))
            
            # call MusicAssembler and start
            ma = MusicAssembler(getRandomBeatFile(filename),app.config['UPLOAD_FOLDER']+'/'+filename)
            # connect to CephFS and register the data
            if (CEPHFS):

            # delete all the temporary data

            # return


if __name__ == "__main__":
    app.run(ssl_context=('cert.pem','key.pem')) # need to run generate_ceert.go here to get cert.pem and key.pem