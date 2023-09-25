const axios = require('axios');
const fs = require('fs');

const PR_NUMBER = process.env.GITHUB_EVENT_PATH ? JSON.parse(fs.readFileSync(process.env.GITHUB_EVENT_PATH, 'utf8')).number : null;
const REPO = process.env.GITHUB_REPOSITORY;
const TOKEN = process.env.GITHUB_TOKEN;

console.log(`https://api.github.com/repos/${REPO}/pulls/${PR_NUMBER}/files`);

axios.get(`https://api.github.com/repos/${REPO}/pulls/${PR_NUMBER}/files`, {})
    .then(response => {
        const files = response.data;
        const root_directories = new Set();

        files.forEach(file => {
            const root_directory = file.filename.split("/")[0];
            root_directories.add(root_directory);
        });

        if (root_directories.size > 1) {
            throw new Error("Multiple root directories changed, action failed.");
        }

        const root_dir = [...root_directories][0];
        console.log(`::set-output name=root_directory::${root_dir}`);
    })
    .catch(error => {
        console.error(error.message || 'Error occurred');
        process.exit(1);
    });
