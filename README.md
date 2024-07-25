<html lang="en">
<body>

<h1>Force Sub Bot</h1>

<a href="https://github.com/Abishnoi69/Force-Sub-Bot/actions?query=workflow%3Abuild+event%3Apush+branch%3Amain"><img src="https://github.com/Abishnoi69/Force-Sub-Bot/workflows/build/badge.svg" alt="build"></a>

<p>This project includes a Telegram bot designed to enforce subscription to a specific channel before allowing users to interact in a group chat. It's built using <a href="https://go.dev">Go</a> and integrates with the Telegram Bot API using <a href="https://github.com/PaulSonOfLars/gotgbot">gotgbot</a>.</p>

<section>
<h2>Installation Instructions</h2>
<h3>Install Go</h3>
<ol>
<li>Clone and install Go:
<pre><code>git clone https://github.com/udhos/update-golang && cd update-golang
sudo ./update-golang.sh
source /etc/profile.d/golang_path.sh</code></pre>
</li>
<li>Exit and reopen your terminal, then verify the installation with <code>go version</code>.</li>
</ol>

<h3>Set Up the Project</h3>
<ol>
<li>Clone the repository:
<pre><code>git clone https://github.com/Abishnoi69/Force-Sub-Bot fSub && cd fSub</code></pre>
</li>
<li>Prepare the environment file:
<ol>
<li>Copy the sample environment file: <code>cp sample.env .env</code></li>
<li>Open env: <code>vi .env</code></li>

<li>Edit the <code>.env</code> file with your preferred editor. Instructions for editing in <code>vi</code>:
<ul>
<li>Press <code>ɪ</code> to start editing.</li>
<li>Press <code>Ctrl + C</code> once editing is complete, then type <code>:wq</code> to save and exit, or <code>:qa</code> to exit without saving.</li>
</ul>
</li>
</ol>
</li>
<li>Start a new <code>tmux</code> session: <code>sudo apt install tmux && tmux</code></li>
<li>Run the bot: <code>go run .</code></li>
</ol>
</section>

<section>
<h2>Deploy to Vercel</h2>
<ol>
<li>Fork this repository 🍴</li>
<li>Login your <a href="https://vercel.com/">Vercel</a> account </li>
<li>Go to your <a href="https://vercel.com/new">Add New Project</a></li>
<li>Choose the repository you forked</li>
<li>Configure the environment variables: <code>DB_URI</code> <a href="https://app.redislabs.com/">Redis</a></li>
<li>Tap on Deploy</li>
</ol>
</section>

<section>
<h2>Usage</h2>
<p>Once the bot is running, it will enforce subscription to a specific channel before allowing users to interact in the group chat. Users not subscribed to the channel will be prompted to do so.</p>

<h3>Commands</h3>
<ul>
<li><code>/start</code> - Start the bot.</li>
<li><code>/fsub</code> - Set the channel for force subscription. Reply to a forwarded message from the channel you wish to set.</li>
<li><code>/fsub on</code> - Enable force subscription mode.</li>
<li><code>/fsub off</code> - Disable force subscription mode.</li>
<li><code>/fsub</code> - Get the current force subscription status.</li>
</ul>
</section>

<section>
<h3>License</h3>
<p>Feel free to visit <a href="LICENSE">LICENSE</a> for more details.</p>
</section>

<section>
<h3>Contributing Guidelines</h3>
<p>Contributions are welcome! For bug reports, feature requests, or pull requests, please open an issue or submit your changes directly.</p>
</section>

<section>
<h3>Support</h3>
<p>For any questions or concerns, please contact me at <a href="https://t.me/Abishnoi1M">Telegram</a>.</p>
<p><a href="https://t.me/FallenAssociation">FallenAssociation</a></p>
</section>

</body>
</html>