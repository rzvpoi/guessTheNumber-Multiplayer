<script setup>
import { onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router';
import { receivedMessage, sendAction } from '../webSocket';
import { useRouter } from 'vue-router';

const router = useRouter();

const messageReceived = ref(false);

let guessedNumbers = ref([])

let guessExists = ref(false)

const route = useRoute();
const code_param = route.params.code;
const code = code_param.split('=')[1]
const player_param = route.params.player
const player = player_param.split('=')[1]
const playerNumber_param = route.params.number
const player_number = playerNumber_param.split('=')[1]

let isPlayerTurn = ref(false)
let latestGuess = 0
let adv_guess = ref(0)

const isDisabled = ref(false);
let game_final_status = ref('')

onMounted(() => {
    var btn = document.getElementById("myBtn");
    var element = document.getElementById("myToast");

    // Create toast instance
    var myToast = new bootstrap.Toast(element);

    btn.addEventListener("click", function () {
        myToast.show();
    });
})

watch(receivedMessage, (newVal) => {
    console.log(window.location.href)
    messageReceived.value = !!newVal; // Update messageReceived based on receivedMessage value
    const msg = receivedMessage._rawValue
    if (msg.event == "run_game") {
        isPlayerTurn.value = true
        isDisabled.value = false
        if (msg.guess_number != null) {
            adv_guess.value = msg.guess_number
            var element = document.getElementById("myToast");
            var myToast = new bootstrap.Toast(element);
            myToast.show();
        }
    } else if (msg.event == "guess_response") {
        guessExists.value = true
        guessedNumbers.value.push({ "number": latestGuess, "guessed": msg.guessed })
        isPlayerTurn.value = false
    } else if (msg.event == "game_finished") {
        console.log("Game finished", msg.isGameWon)
        if (msg.isGameWon == true) {
            game_final_status.value = "You won!"
        } else {
            game_final_status.value = "You lost :("
        }
        document.getElementById('game_status_btn').click()
    }
});

function returnToMenu() {
    router.push('/')
}

let numberGuess = ref('')

function sendValue() {
    guessNumber(numberGuess.value)
}

function guessNumber(number) {
    if (number.length != 4) {
        document.getElementById('errorGuess').innerText = "Number must have 4 digits"
        return
    }
    if(isNaN(number) ){
        document.getElementById('errorGuess').innerText = "Value must be a number"
        return
    } 
    isDisabled.value = true
    document.getElementById('errorGuess').innerText = ""
    document.getElementById('guessInput').value = ''
    latestGuess = number
    sendAction({ action: 'guess_number', code: code, player: player, guess: number })
}

</script>

<template>
    <button style="display: none;" type="button" class="btn btn-primary" id="myBtn">Show Toast</button>
    <div class="toast" id="myToast">
        <div class="toast-header">
            <strong class="me-auto"><i class="bi-gift-fill"></i>Opponent number guess</strong>
            <small></small>
            <button type="button" class="btn-close" data-bs-dismiss="toast"></button>
        </div>
        <div class="toast-body">
            Your opponent tried number {{adv_guess}}
        </div>
    </div>


    <!-- Button trigger modal -->
    <button type="button" style="display: none;" class="btn btn-primary" data-bs-toggle="modal"
        data-bs-target="#staticBackdrop" id="game_status_btn">
    </button>

    <!-- Modal -->
    <div class="modal fade" id="staticBackdrop" data-bs-backdrop="static" data-bs-keyboard="false" tabindex="-1"
        aria-labelledby="staticBackdropLabel" aria-hidden="true">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-header">
                    <h1 class="modal-title fs-5" id="staticBackdropLabel">{{ game_final_status }}</h1>
                    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                </div>
                <div class="modal-body">
                    Stats:
                </div>
                <div class="modal-footer">
                    <button @click="returnToMenu" type="button" class="btn btn-primary" data-bs-dismiss="modal">Return to
                        main menu</button>
                </div>
            </div>
        </div>
    </div>


    <div class="container">
        <!-- Title at the top left -->
        <div class="row mt-3">
            <div class="col">
                <h1>Your number is {{ player_number }}</h1>
            </div>
        </div>

        <!-- Row for main content -->
        <div class="row mt-1">
            <!-- Container for list of numbers -->
            <div class="col-md-6">
                <div class="card" style="max-height: 600px; overflow-y: auto;">
                    <div class="card-body">
                        <h5 class="card-title">
                            <div v-if="!isPlayerTurn">
                                <input type="text" class="form-control" placeholder="Waiting for other player to guess"
                                    aria-label="Waiting for other player to guess" readonly>
                            </div>
                            <div v-else>
                                <div class="input-group mb-0">
                                    <input v-model="numberGuess" type="text" class="form-control"
                                        placeholder="Guess the number" aria-label="Guess the number"
                                        aria-describedby="button-addon2" id="guessInput">
                                    <button @click="sendValue" class="btn btn-outline-secondary" type="button"
                                        id="button-addon2" :disabled="isDisabled">Guess</button>
                                </div>
                                <small id="errorGuess" class="form-text"
                                    style="color: red;font-size: 0.8rem; overflow-y: auto;"></small>
                            </div>

                        </h5>
                        <p>Your guesses:</p>
                        <ul v-if="guessExists" class="list-group">
                            <li v-for="number in guessedNumbers" class="list-group-item">{{ number.number }} - {{
                                number.guessed }} number guessed</li>
                        </ul>
                        <ul v-else>
                            <p>No numbers for now...</p>
                        </ul>
                    </div>
                </div>
            </div>

            <!-- Container to the right side -->
            <div class="col-md-6">
                <div class="card" style="max-height: 300px; overflow-y: auto;">
                    <div class="card-body">
                        <h5 class="card-title">Notes</h5>
                        <div class="input-group">
                            <textarea class="form-control" aria-label="With textarea"></textarea>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
  
  
<style scoped>

.toast{
    position: absolute; 
    top: 10px; 
    right: 10px;
}

</style>
  