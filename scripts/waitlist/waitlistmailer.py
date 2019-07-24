import sendgrid
import os
import base64

from sendgrid.helpers.mail import *

sg = sendgrid.SendGridAPIClient(apikey=os.environ.get('SENDGRID_API_KEY'))
from_email = Email("team@october.news")

users = [
    {'email': 'konrad@october.news'}
]

file_path = 'october.jpg'
with open(file_path,'rb') as f:
    data = f.read()
    f.close()

encoded = base64.b64encode(data).decode()

attachment = Attachment()
attachment.content = encoded
attachment.type = "image/jpeg"
attachment.filename = "october.jpg"
attachment.disposition = "inline"
attachment.content_id = "banner"

for user in users:
    to_email = Email(user['email'])
    subject = "Your Invitation to October"
    content = Content("text/html", """<img src="cid:banner"><br><br>
                      Dear Konrad,<br><br>

We're ready to add early testers to October off the waitlist, and it's your turn.<br><br>

We are giving you a personal invite code. Go to https://october.app/qr/9EF57 to get started. iPhone preferred, but Desktop Web also works.<br><br>

October is a social network where it's easy to share secrets and have difficult conversations using a mixture of real names and pseudonyms.<br><br>

Good content earns you coins that you can use to post without revealing your name. Think of it as a safe place to discuss difficult subjects.<br><br>

Our goal is to foster communication across tribal boundaries and build new bridges between partisan communities.<br><br>

Feel free to get in touch with us if anything breaks. We're honored to have you as an early tester!<br><br>

Team October""")
    mail = Mail(from_email, subject, to_email, content)
    mail.add_attachment(attachment)
    response = sg.client.mail.send.post(request_body=mail.get())
    print(response.status_code)
    print(response.body)
    print(response.headers)
