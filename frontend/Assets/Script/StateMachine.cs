using System;
using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using static GoogleCloudStreamingSpeechToText.StreamingRecord;


public class StateMachine : MonoBehaviour
{
    public AudioSource audioSource;
    private enum gameState
    {
        standby,
        recognition,
        awaitReply,
        reply
    }
    private gameState _currentState = gameState.standby;


    private void LateUpdate()
    {
        switch(_currentState)
        {
            case gameState.standby://not user 
                
                
                break;
            case gameState.recognition://user talk start

                break;
            case gameState.awaitReply://user talk end and chatgpt

                break;
            case gameState.reply://avater talk

                break;
        }
    }
    public async void EndRecognition(string finalDetection)
    {
        _currentState = gameState.awaitReply;
        var text = await HttpTest.GuideByText(finalDetection);
        Debug.Log(text);
        var test = await HttpTest.TextToAudio("1", "true", text);
        byte[] binaryData = Convert.FromBase64String(test);
        AudioClip clip = Wav.ToAudioClip(binaryData, "test");
        audioSource.PlayOneShot(clip);
    }
    public void StartRecognition()
    {
        _currentState= gameState.recognition;

    }


}